<script>
    import { page } from '$app/stores';

    // Access the 'id' from the URL: /basket/[id]
    $: basketId = $page.params.id;
    $: shareUrl = typeof window !== 'undefined' ? window.location.href : '';
</script>

<h1>Your Basket</h1>
<p>Unique ID: <code>{basketId}</code></p>

<div class="share-box">
    <p>Share this link to collaborate:</p>
    <input type="text" readonly value={shareUrl} />
    <button on:click={() => navigator.clipboard.writeText(shareUrl)}>
        Copy Link
    </button>
</div>

<!-- Your shopping list logic goes here -->