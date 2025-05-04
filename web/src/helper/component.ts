import { KeyboardEvent } from 'react';

// Stop the Enter key from submitting the form
export const stopEnterPropagation = (event: KeyboardEvent<HTMLInputElement>): void => {
    if (event.key === 'Enter') {
        event.preventDefault(); // Prevent the default action (form submission/reload)
        // Optionally, you could trigger a manual filter submission here if needed,
        // but FilterLiveSearch usually handles it live.
        // For example: event.target.blur(); // Remove focus from input
    }
};