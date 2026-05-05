<script>
    import { goto } from '$app/navigation';

    function generateUUID() {
        // crypto is not available when not run in Secure Context (SSL)
        if (typeof crypto !== 'undefined' && crypto.randomUUID) {
            return crypto.randomUUID();
        }
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
            const r = Math.random() * 16 | 0;
            const v = c === 'x' ? r : (r & 0x3 | 0x8);
            return v.toString(16);
        });
    }

    function createNewBasket() {
        const uniqueId = generateUUID();
        // Redirect the user to the new basket URL
        goto(`/basket/${uniqueId}`);
    }
</script>

<div class="landing-container">
    <div class="main-content">
        <h1 class="title">Shopping List</h1>
        <p class="subtitle">Create a private basket and share the URL with others to collaborate in real-time.</p>

        <button class="btn-create" on:click={createNewBasket}>
            Create New Basket
        </button>
    </div>

    <a
            class="source-link"
            href="https://github.com/ubap/go-shop"
            target="_blank"
            rel="noopener noreferrer"
    >
        <svg viewBox="0 0 24 24" width="20" height="20" stroke="currentColor" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round" class="github-icon">
            <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"></path>
        </svg>
        View Source Code
    </a>
</div>

<style>
    .landing-container {
        max-width: 600px;
        margin: 0 auto;
        padding: 40px 20px;
        font-family: system-ui, -apple-system, sans-serif;
        box-sizing: border-box;
        text-align: center;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: space-between;
        min-height: 80vh;
    }

    .main-content {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        flex-grow: 1;
    }

    .title {
        font-size: 2.5rem;
        margin-bottom: 16px;
        color: #333;
        font-weight: 800;
        margin-top: 0;
    }

    .subtitle {
        font-size: 1.1rem;
        color: #666;
        margin-bottom: 40px;
        line-height: 1.5;
        max-width: 400px;
    }

    .btn-create {
        padding: 16px 32px;
        font-size: 1.2rem;
        font-weight: 600;
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 8px;
        cursor: pointer;
        transition: background-color 0.2s, transform 0.1s;
        box-shadow: 0 4px 12px rgba(0, 123, 255, 0.2);
    }

    .btn-create:hover {
        background-color: #0056b3;
    }

    .btn-create:active {
        transform: translateY(2px);
    }

    .source-link {
        display: flex;
        align-items: center;
        gap: 8px;
        color: #888;
        text-decoration: none;
        font-size: 0.95rem;
        font-weight: 500;
        padding: 12px;
        border-radius: 6px;
        transition: color 0.2s, background-color 0.2s;
        margin-top: 40px;
    }

    .source-link:hover {
        color: #333;
        background-color: #f5f5f5;
    }
</style>