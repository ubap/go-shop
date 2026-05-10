<script lang="ts">
    import { page } from '$app/stores';
    import { dev } from '$app/environment';
    import { afterNavigate } from '$app/navigation';
    import { flip } from 'svelte/animate';
    import { slide, fade, fly } from 'svelte/transition';

    interface Item {
        id: number;
        name: string;
        completed: boolean;
    }

    let items: Item[] = $state([]);
    let newItem: string = $state("");
    let inputRef: HTMLInputElement | undefined = $state();

    let socket: WebSocket | undefined;
    let isConnected = $state(false);

    let showUndo = $state(false);
    let lastDeletedItem = $state<Item | null>(null);
    let undoTimeout : ReturnType<typeof setTimeout> | undefined;

    $effect(() => {
        console.log(window.__INITIAL_DATA__)
        if (window.__INITIAL_DATA__) {
            items = window.__INITIAL_DATA__;
        }

        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        let wsUrl = "";
        if (dev) {
            const hostname = window.location.hostname;
            wsUrl = `${protocol}//${hostname}:9090/ws?id=${$page.params.id}`;
        } else {
            const host = window.location.host;
            wsUrl = `${protocol}//${host}/ws?id=${$page.params.id}`;
        }

        let isMounted = true;
        let reconnectTimeout: ReturnType<typeof setTimeout>;

        function connect() {
            if (!isMounted || isSocketReady()) return;
            console.log(`Connecting to ws ...`);
            socket = new WebSocket(wsUrl);
            socket.onopen = () => {
                isConnected = true;
            };
            socket.onmessage = (event) => {
                const data = JSON.parse(event.data);
                if (data.type === "full_list") {
                    items = data.items;
                }
            };
            socket.onclose = () => {
                isConnected = false;
                if (isMounted) {
                    console.log("WebSocket disconnected. Reconnecting in 2 seconds...");
                    clearTimeout(reconnectTimeout);
                    reconnectTimeout = setTimeout(connect, 2000);
                }
            };
            socket.onerror = () => {
                socket?.close();
            };
        }
        document.addEventListener('visibilitychange', () => {
            // connect instantly after the page is unminimized on mobile
            if (document.visibilityState === 'visible') {
                connect();
            }
        });
        connect();

        return () => {
            isMounted = false;
            clearTimeout(reconnectTimeout);
            if (socket) socket.close();
        };
    });

    afterNavigate(() => {
        inputRef?.focus();
    });

    function isSocketReady() {
        return socket && socket.readyState === WebSocket.OPEN;
    }

    function addItem() {
        if (!newItem || newItem.trim().length === 0 || !isSocketReady()) return;

        socket!.send(JSON.stringify({
            type: "addItem",
            itemName: newItem
        }));

        newItem = "";
        inputRef?.focus();
    }

    function toggleItem(itemId: number, currentState: boolean) {
        if (!isSocketReady()) return;

        socket!.send(JSON.stringify({
            type: "setItemCompletion",
            id: itemId,
            completed: !currentState
        }));
    }

    function removeItem(item: Item) {
        if (!isSocketReady()) return;

        lastDeletedItem = item;
        showUndo = true;
        clearTimeout(undoTimeout);
        undoTimeout = setTimeout(() => {
            showUndo = false;
        }, 5000);

        socket!.send(JSON.stringify({
            type: "deleteItem",
            id: item.id
        }));
    }

    function restoreItem() {
        if (lastDeletedItem == null) return;
        socket!.send(JSON.stringify({
            type: "restoreItem",
            id: lastDeletedItem.id
        }));
        showUndo = false;
        clearTimeout(undoTimeout);

    }
</script>

<header class="app-header">
    <div class="header-content">
    <div class="header-brand">
        <span class="logo">🛒</span>
        <div>
            <h1 class="app-title">Shopping List</h1>
        </div>
    </div>

    <div class="header-user">
    </div>
    </div>
</header>

<div class="shoplist-container">
    {#if !isConnected}
        <div class="connection-overlay" transition:fade={{ duration: 200 }}>
            <div class="loader"></div>
            <p class="loader-text">Connecting...</p>
        </div>
    {/if}

    <div class="input-group">
        <input
                class="item-input"
                type="search"
                autocomplete="off"
                enterkeyhint="done"

                bind:this={inputRef}
                bind:value={newItem}
                placeholder="Add item..."
                on:keydown={(e) => e.key === 'Enter' && addItem()}
                disabled={!isConnected}
        />
        <button class="btn-add" on:click={addItem} disabled={!isConnected}>Add</button>
    </div>

    <ul class="item-list">
        {#each items as item (item.id)}
            <li class={item.completed ? 'completed' : ''}
                animate:flip={{ duration: 300 }}
                transition:slide|local>
                <label class="row-label">
                    <input
                            class="checkbox"
                            type="checkbox"
                            checked={item.completed}
                            on:change={() => toggleItem(item.id, item.completed)}
                            disabled={!isConnected}
                    />
                    <span class="title">{item.name}</span>
                </label>

                <button class="btn-delete" on:click={() => removeItem(item)} disabled={!isConnected}>✕</button>
            </li>
        {/each}
    </ul>

    <div class="list-footer">
        <span class="footer-id">List ID: {$page.params.id}</span>
        <br/>
        <button class="btn-copy" on:click={() => {
            navigator.clipboard.writeText(window.location.href);
        }}>
            Copy Share URL
        </button>
    </div>
</div>

{#if showUndo}
    <div class="undo-toast" transition:fly={{ y: 50, duration: 300 }}>
        <span class="undo-text">Deleted: <strong>{lastDeletedItem?.name}</strong></span>
        <button class="undo-button" on:click={restoreItem}>
            Undo
        </button>
    </div>
{/if}

<style>
    .shoplist-container {
        max-width: 600px;
        margin: 0 auto;
        padding: 20px;
        font-family: system-ui, -apple-system, sans-serif;
        box-sizing: border-box;
        position: relative;
    }

    .connection-overlay {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(255, 255, 255, 0.75);
        backdrop-filter: blur(2px);
        z-index: 10;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        border-radius: 8px;
    }

    .loader {
        border: 4px solid #f3f3f3;
        border-top: 4px solid #007bff;
        border-radius: 50%;
        width: 36px;
        height: 36px;
        animation: spin 1s linear infinite;
        margin-bottom: 12px;
    }

    .loader-text {
        color: #007bff;
        font-weight: 600;
        margin: 0;
    }

    @keyframes spin {
        0% { transform: rotate(0deg); }
        100% { transform: rotate(360deg); }
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
        flex-shrink: 0;
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

    .btn-add:hover:not(:disabled) {
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

    .btn-delete:hover:not(:disabled) {
        background-color: #ffe6e6;
    }

    input:disabled, button:disabled {
        cursor: not-allowed;
    }

    .list-footer {
        margin-top: 32px;
        padding-top: 16px;
        border-top: 1px dashed #ddd;
        display: flex;
        justify-content: space-between;
        align-items: center;
        font-size: 12px;
        color: #888;
    }

    .footer-id {
        font-family: monospace;
    }

    .btn-copy {
        background: none;
        border: 1px solid #ccc;
        padding: 6px 12px;
        border-radius: 4px;
        font-size: 12px;
        color: #555;
        cursor: pointer;
        transition: all 0.2s;
    }

    .btn-copy:hover {
        background: #e9ecef;
        color: #333;
    }

    .header-content {
        max-width: 600px;
        margin: 0 auto;

        padding: 6px 20px;

        box-sizing: border-box;
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;

        height: 32px;
    }

    .header-user {
        display: flex;
        align-items: center;
        gap: 8px;
        cursor: pointer;
    }

    .app-header {
        font-family: system-ui, -apple-system, sans-serif;

        border-bottom: 1px solid #eaeaea;
        background-color: #ffffff;
        padding: 12px 20px;
    }

    .header-brand {
        display: flex;
        align-items: center;
        gap: 12px;
    }

    .logo {
        font-size: 20px;
        line-height: 1;
        display: flex;
        align-items: center;
    }

    .app-title {
        margin: 0 !important;
        font-size: 1.15rem;
        font-weight: 600;
        line-height: 1;
        color: #333;
    }

    .undo-toast {
        position: fixed;
        bottom: 20px; /* Odstęp od dołu ekranu */
        left: 50%;
        transform: translateX(-50%); /* Centrowanie w poziomie */

        background-color: #333; /* Ciemne tło, żeby się odróżniało */
        color: white;
        padding: 12px 20px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        gap: 16px;
        z-index: 9999; /* Musi być nad wszystkim innym */
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        min-width: 280px;
        justify-content: space-between;

        font-family: system-ui, -apple-system, sans-serif;
    }

    .undo-text {
        font-size: 14px;
    }

    .undo-button {
        background: none;
        border: none;
        color: #4dabf7; /* Jasnoniebieski kolor akcji */
        font-weight: 800;
        cursor: pointer;
        padding: 4px 8px;
        font-size: 13px;
        letter-spacing: 0.5px;
    }

    .undo-button:hover {
        background-color: rgba(255,255,255,0.1);
        border-radius: 4px;
    }
</style>