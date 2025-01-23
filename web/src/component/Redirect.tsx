import { useEffect } from "react";
import { useRedirect } from "react-admin";

export default function Redirect({ to }: { to: string }) {
    const redirect = useRedirect()

    useEffect(
        () => {
            redirect(to)
        }, [to]
    )

    return (
        <>
        </>
    )

}