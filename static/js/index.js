async function uploadFiles(fileInput) {
    const data = new FormData();
    for (const file of fileInput.files) {
        data.append('messenger_files', file, file.name);
    }

    const r = await fetch('/api/upload', {
        method: 'POST',
        body: data,
    });
    const json = await r.json();

    if (json.error !== '') {
        alert(`Error: ${json.error}`);
        return;
    }

    window.location.href = `/stats?id=${json.id}`;
}

window.addEventListener('load', async () => {
    const fileInput = document.getElementById('messenger-file-input');
    const uploadBtn = document.getElementById('get-stats-btn');

    uploadBtn.addEventListener('click', async () => { await uploadFiles(fileInput); });
});
