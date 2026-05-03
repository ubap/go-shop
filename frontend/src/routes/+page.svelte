<script>
    import { goto } from '$app/navigation';

    function createNewBasket() {
        // crypto is not available when not run in Secure Context (SSL)
        const uniqueId = crypto.randomUUID ? crypto.randomUUID() : Math.random().toString(36).substring(2);

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