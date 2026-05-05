<script lang="ts">
    import { page } from '$app/stores';
    import { dev } from '$app/environment';
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
        let wsUrl = "";
        if (dev) {
            const hostname = window.location.hostname;
            wsUrl = `${protocol}//${hostname}:9090/ws?id=${$page.params.id}`;
        } else {
            const host = window.location.host;
            wsUrl = `${protocol}//${host}/ws?id=${$page.params.id}`;
        }

        socket = new WebSocket(wsUrl);
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

<div class="shoplist-container">
    <h1 class="header-title">Basket: {$page.params.id}</h1>

    <div class="input-group">
        <input
                class="item-input"
                bind:this={inputRef}
                bind:value={newItem}
                placeholder="Add item..."
                on:keydown={(e) => e.key === 'Enter' && addItem()}
        />
        <button class="btn-add" on:click={addItem}>Add</button>
    </div>

    <ul class="item-list">
        {#each items as item (item.id)}
            <li class={item.completed ? 'completed' : ''}>
                <label class="row-label">
                    <input
                            class="checkbox"
                            type="checkbox"
                            checked={item.completed}
                            on:change={() => toggleItem(item.id, item.completed)}
                    />
                    <span class="title">{item.name}</span>
                </label>

                <button class="btn-delete" on:click={() => removeItem(item.id)}>✕</button>
            </li>
        {/each}
    </ul>
</div>

<style>
    .shoplist-container {
        max-width: 600px;
        margin: 0 auto;
        padding: 20px;
        font-family: system-ui, -apple-system, sans-serif;
        box-sizing: border-box;
    }

    .header-title {
        font-size: 1.5rem;
        margin-bottom: 24px;
        color: #333;
    }

    .input-group {
        display: flex;
        gap: 10px;
        margin-bottom: 24px;
        width: 100%;
    }

    .item-input {
        flex-grow: 1;
        min-width: 0;
        box-sizing: border-box;
        padding: 12px 16px;
        font-size: 16px;
        border: 1px solid #ccc;
        border-radius: 8px;
        outline: none;
        transition: border-color 0.2s;
    }

    .item-input:focus {
        border-color: #007bff;
    }

    .btn-add {
        flex-shrink: 0; /* KLUCZOWE: Zapobiega zgnieceniu przycisku, gdy brakuje miejsca */
        box-sizing: border-box;
        padding: 12px 24px;
        font-size: 16px;
        font-weight: 600;
        background-color: #007bff;
        color: white;
        border: none;
        border-radius: 8px;
        cursor: pointer;
        transition: background-color 0.2s;
    }

    .btn-add:hover {
        background-color: #0056b3;
    }

    .item-list {
        list-style: none;
        padding: 0;
        margin: 0;
        display: flex;
        flex-direction: column;
        gap: 8px;
    }

    li {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 12px 16px;
        border-radius: 8px;
        background-color: #f9f9f9;
        border: 1px solid #eaeaea;
        transition: background-color 0.2s, opacity 0.2s;
    }

    .row-label {
        display: flex;
        align-items: center;
        gap: 12px;
        cursor: pointer;
        flex-grow: 1;
        user-select: none;
    }

    .checkbox {
        width: 20px;
        height: 20px;
        cursor: pointer;
        margin: 0;
    }

    .title {
        font-size: 16px;
        color: #333;
        transition: color 0.2s;
    }

    .completed {
        background-color: #f1f1f1;
        border-color: #f1f1f1;
        opacity: 0.7;
    }

    .completed .title {
        text-decoration: line-through;
        color: #888;
    }

    .btn-delete {
        background: transparent;
        border: none;
        color: #ff4d4f;
        font-size: 18px;
        cursor: pointer;
        padding: 8px;
        border-radius: 4px;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: background-color 0.2s;
    }

    .btn-delete:hover {
        background-color: #ffe6e6;
    }
</style>