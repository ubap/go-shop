<script lang="ts">
    import { page } from '$app/stores';
    import { afterNavigate } from '$app/navigation';

    interface Item {
        id: number;
        name: string;
        completed: boolean;
    }

    let items: Item[] = $state([]);
    let newItem: string = $state("");
    let inputRef: HTMLInputElement | undefined = $state();

    let socket: WebSocket | undefined; 

    $effect(() => {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const host = window.location.hostname;
        const wsPort = "8080";

        socket = new WebSocket(`${protocol}//${host}:${wsPort}/ws?id=${$page.params.id}`);

        socket.onmessage = (event) => {
            const data = JSON.parse(event.data);
            if (data.type === "full_list") {
                items = data.items;
            }
        };

        return () => {
            // cleanup
            if (socket) socket.close();
        };
    });

    afterNavigate(() => {
        inputRef?.focus();
    });

    function addItem() {
        if (!newItem || newItem.trim().length === 0 || !socket) return;

        socket.send(JSON.stringify({
            type: "addItem",
            itemName: newItem
        }));

        newItem = "";
        inputRef?.focus();
    }

    function toggleItem(itemId: number, currentState: boolean) {
        if (!socket) return;
        socket.send(JSON.stringify({
            type: "setItemCompletion",
            id: itemId,
            completed: !currentState
        }));
    }

    function removeItem(itemId: number) {
        if (!socket) return;
        socket.send(JSON.stringify({
            type: "deleteItem",
            id: itemId
        }));
    }
</script>

<h1>Basket: {$page.params.id}</h1>

<div class="input-group">
    <input 
        bind:this={inputRef}
        bind:value={newItem} 
        placeholder="Add item..." 
        on:keydown={(e) => e.key === 'Enter' && addItem()} 
    />
    <button on:click={addItem}>Add</button>
</div>

<ul>
    {#each items as item (item.id)}
        <li class={item.completed ? 'completed' : ''}>
            <label class="row-label">
                <input
                    type="checkbox"
                    checked={item.completed}
                    on:change={() => toggleItem(item.id, item.completed)}
                />
                <span class="title">{item.name}</span>
            </label>

            <button on:click={() => removeItem(item.id)}>x</button>
        </li>
    {/each}
</ul>

<style>
    .input-group {
        display: flex;
        gap: 10px;
        margin-bottom: 20px;
    }
    
    .completed .title {
        text-decoration: line-through;
        color: gray;
    }
    
    li {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 8px;
        padding: 5px;
        border-radius: 4px;
        background-color: #fcfcfc;
    }
    
    .row-label {
        display: flex;
        align-items: center;
        gap: 10px;
        cursor: pointer;
        flex-grow: 1;
        padding: 5px 0;
    }
    
    button {
        cursor: pointer;
    }
</style>