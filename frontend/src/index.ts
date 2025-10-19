document.addEventListener('DOMContentLoaded', function () {
    const itemInput = document.querySelector<HTMLAnchorElement>('#new-basket-btn')!;
    itemInput.href = "/basket.html?id=" + crypto.randomUUID()
});
