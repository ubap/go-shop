<script>
    import { goto } from '$app/navigation';

    function createNewBasket() {
        // Generate a unique ID (e.g., "550e8400-e29b-41d4-a716-446655440000")
        const uniqueId = crypto.randomUUID();

        // Redirect the user to the new basket URL
        goto(`/basket/${uniqueId}`);
    }
</script>

<h1>Welcome to ShoppingList</h1>
<p>Create a private basket and share the URL with others.</p>

<button on:click={createNewBasket}>
    Create New Basket
</button>

<style>
    button {
        padding: 10px 20px;
        font-size: 1.2rem;
        cursor: pointer;
        background: #ff3e00;
        color: white;
        border: none;
        border-radius: 4px;
    }
</style>