import {
    Button,
    useCreate,
    useGetManyReference,
    useListContext,
    useNotify,
    useRefresh,
    useUnselectAll
} from "react-admin";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import KeyBoardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import React from "react";
import {Button} from "@mui/material";

export default function SendSpecimenToWorkOrder() {
    const {selectedIds} = useListContext();
    const refresh = useRefresh();
    const notify = useNotify();
    const unselectAll = useUnselectAll('Specimen');
    const {data} = useGetManyReference<any>(
        'work-order',
        {
            target: 'Specimen_ids',
            id: selectedIds.join(","),
            pagination: {page: 1, perPage: 100},
            filter: {status: 'pending'},
            sort: {field: 'created_at', order: 'DESC'}
        }
    );
    const [create, {isPending}] = useCreate();
    const sendToOrder = (order: any) => {
        return async () => {
            await create(
                `work-order/${order.id}/Specimen`,
                {
                    data: {
                        Specimen_ids: selectedIds
                    },
                },
                {
                    onSuccess: () => {
                        notify('Success send to order');
                        unselectAll();
                    },
                    onError: () => {
                        notify('Failed to send to order', {type: 'error'});
                        refresh();
                    },
                }
            );
        }
    }

    const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
    const open = Boolean(anchorEl);
    const handleClick = (event: React.MouseEvent<HTMLElement>) => {
        setAnchorEl(event.currentTarget);
    };
    const handleClose = () => {
        setAnchorEl(null);
    };

    return (
        <div>
            <Button
                id="demo-customized-button"
                aria-controls={open ? 'send-to-order' : undefined}
                aria-haspopup="true"
                aria-expanded={open ? 'true' : undefined}
                variant="contained"
                disableElevation
                onClick={handleClick}
                endIcon={<KeyBoardArrowDownIcon/>}
                label={"Send to order"}
                disabled={isPending}
            ></Button>
            <Menu
                id="demo-customized-menu"
                MenuListProps={{
                    'aria-labelledby': 'send-to-order',
                }}
                anchorEl={anchorEl}
                open={open}
                onClose={handleClose}>
                {data?.map((order: any) => (
                    <MenuItem key={order.id} onClick={sendToOrder(order)}>
                        #{order.id}-{order.description}
                    </MenuItem>
                ))}
            </Menu>
        </div>
    );
};
