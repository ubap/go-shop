const app = document.querySelector<HTMLDivElement>('#app')!;
app.innerHTML = ` 
	<h1>Hello from TypeScript!</h1> 
`;


// Select the element once and reuse it. It's more efficient.
const apiResultDiv = document.querySelector<HTMLDivElement>('#api-call-result');

async function apiCall(): Promise<void> {
    // Make sure the element was actually found before trying to use it.
    if (!apiResultDiv) {
        console.error("Could not find the element with ID #api-call-result");
        return;
    }

    try {
        // "await" pauses the function until the fetch is complete.
        const response = await fetch("/api");

        // Check if the network request was successful.
        if (!response.ok) {
            // Throw an error to be caught by the catch block.
            throw new Error(`Network response was not ok: ${response.statusText}`);
        }

        // "await" pauses again until the response body is read as text.
        const message = await response.text();

        // Now, update the DOM with the final text content.
        apiResultDiv.innerHTML = message;

    } catch (error) {
        // Handle any errors that occurred during the fetch.
        console.error("There was a problem with the fetch operation:", error);
        apiResultDiv.innerHTML = "Failed to fetch data from the API.";
    }
}

// Example of how to call it:
apiCall();