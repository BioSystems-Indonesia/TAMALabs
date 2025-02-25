import * as React from 'react';
import Box from '@mui/material/Box';
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepLabel from '@mui/material/StepLabel';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import ArrowForwardIosIcon from '@mui/icons-material/ArrowForwardIos';
import ArrowBackIosIcon from '@mui/icons-material/ArrowBackIos';
import { useFormContext } from "react-hook-form";


export type FormStepperProps = {
    activeStep: number;
    setActiveStep: React.Dispatch<React.SetStateAction<number>>;
    children: React.ReactNode;
    steps: string[];
    onFinish?: (data: any) => void
};

export default function FormStepper({
    activeStep,
    setActiveStep,
    children,
    steps,
    onFinish,
}: FormStepperProps) {
    const isStepOptional = (step: number) => {
        return false;
    };

    const handleNext = () => {
        setActiveStep((prevActiveStep) => prevActiveStep + 1);
    };

    const handleBack = () => {
        setActiveStep((prevActiveStep) => prevActiveStep - 1);
    };

    const { getValues } = useFormContext();

    return (
        <Box sx={{ width: '100%' }}>
            <Stepper activeStep={activeStep}>
                {steps.map((label, index) => {
                    const stepProps: { completed?: boolean } = {};
                    const labelProps: {
                        optional?: React.ReactNode;
                    } = {};
                    if (isStepOptional(index)) {
                        labelProps.optional = (
                            <Typography variant="caption">Optional</Typography>
                        );
                    }
                    return (
                        <Step key={label} {...stepProps}>
                            <StepLabel {...labelProps}>{label}</StepLabel>
                        </Step>
                    );
                })}
            </Stepper>
            <React.Fragment>
                <Typography sx={{ mt: 2, mb: 1 }}>Step {activeStep + 1}</Typography>
                {children}
                <Box sx={{ display: 'flex', flexDirection: 'row', pt: 2 }}>
                    <Button
                        color="inherit"
                        disabled={activeStep === 0}
                        onClick={handleBack}
                        sx={{ mr: 1, width: '120px' }}
                        variant='contained'
                        startIcon={<ArrowBackIosIcon />}
                    >
                        Back
                    </Button>
                    <Box sx={{ flex: '1 1 auto' }} />
                    <Button onClick={
                        activeStep === steps.length - 1 ? () => {
                            if (onFinish) {
                                const { ...data } = getValues();
                                onFinish(data)
                            }
                        } :
                            handleNext
                    } variant='contained'
                        sx={{
                            width: '120px',
                        }}
                        endIcon={<ArrowForwardIosIcon />}
                    >
                        {activeStep === steps.length - 1 ? 'Finish' : 'Next'}
                    </Button>
                </Box>
            </React.Fragment>
        </Box>
    );
}
