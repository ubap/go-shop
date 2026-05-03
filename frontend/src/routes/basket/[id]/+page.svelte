<script lang="ts">
    import { onMount, onDestroy } from 'svelte';
    import { page } from '$app/stores';

    interface Item {
        id: number;        // maps to Go int64
        name: string;      // maps to Go Title (json:"name")
        completed: boolean; // maps to Go Completed
    }

    let id = $page.params.id;
    let socket: WebSocket;
    let items: Item[] = [];
    let newItem = "";

    onMount(() => {
        socket = new WebSocket(`ws://localhost:8080/ws?id=${id}`);

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.type === "full_list") {
                items = data.items; // UI updates automatically!
            }
        };
    });

    onDestroy(() => {
        if (socket) socket.close();
    });

    function addItem() {
        if (!newItem || newItem.trim().length === 0) return;

        socket.send(JSON.stringify({
            type: "addItem",
            itemName: newItem
        }));

        newItem = "";
    }

    function toggleItem(itemId: number, currentState: boolean) {
        socket.send(JSON.stringify({
            type: "setItemCompletion",
            id: itemId,
            completed: !currentState
        }));
    }

    function removeItem(itemId: number) {
        socket.send(JSON.stringify({
            type: "deleteItem",
            id: itemId
        }));
    }

    function focusOnInit(node: HTMLInputElement) {
        node.focus();
    }
</script>

<h1>Basket: {id}</h1>

<div class="input-group">
    <input use:focusOnInit bind:value={newItem} placeholder="Add item..." on:keydown={(e) => e.key === 'Enter' && addItem()} />
    <button on:click={addItem}>Add</button>
</div>

<ul>
    {#each items as item (item.id)}
        <li class={item.completed ? 'completed' : ''}>
            <input
                    type="checkbox"
                    checked={item.completed}
                    on:change={() => toggleItem(item.id, item.completed)}
            />

            <span class="title">{item.name}</span>

            <button on:click={() => removeItem(item.id)}>x</button>
        </li>
    {/each}
</ul>

<style>
    .completed .title {
        text-decoration: line-through;
        color: gray;
    }
    li {
        display: flex;
        align-items: center;
        gap: 10px;
        margin-bottom: 5px;
    }
</style>