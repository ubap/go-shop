<script>
    import { onMount, onDestroy } from 'svelte';
    import { page } from '$app/stores';

    let id = $page.params.id;
    let socket;
    let items = [];
    let newItem = "";

    onMount(() => {
        // Connect to Go WebSocket
        socket = new WebSocket(`ws://localhost:8080/ws?id=${id}`);

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.type === "update") {
                items = data.items; // UI updates automatically!
            }
        };
    });

    onDestroy(() => {
        if (socket) socket.close();
    });

    function addItem() {
        if (!newItem) return;
        const updatedItems = [...items, newItem];

        // Send the updated list to Go
        socket.send(JSON.stringify({
            type: "update",
            items: updatedItems
        }));

        newItem = "";
    }

    function removeItem(index) {
        const updatedItems = items.filter((_, i) => i !== index);
        socket.send(JSON.stringify({ type: "update", items: updatedItems }));
    }
</script>

<h1>Basket: {id}</h1>
<p>Share this URL. Anyone with the link can edit live!</p>

<input bind:value={newItem} placeholder="Add item..." />
<button on:click={addItem}>Add</button>

<ul>
    {#each items as item, i}
        <li>
            {item}
            <button on:click={() => removeItem(i)}>x</button>
        </li>
    {/each}
</ul>