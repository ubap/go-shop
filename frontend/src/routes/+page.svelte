<script>
    import { goto } from '$app/navigation';

    function createNewBasket() {
        // crypto is not available when not run in Secure Context (SSL)
        const uniqueId = crypto.randomUUID ? crypto.randomUUID() : Math.random().toString(36).substring(2);

        // Redirect the user to the new basket URL
        goto(`/basket/${uniqueId}`);
    }
</script>

<div class="landing-container">
    <h1 class="title">Shopping List</h1>
    <p class="subtitle">Create a private basket and share the URL with others to collaborate in real-time.</p>

    <button class="btn-create" on:click={createNewBasket}>
        Create New Basket
    </button>
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
        justify-content: center;
        min-height: 80vh;
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
</style>