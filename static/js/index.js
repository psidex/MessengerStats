const fileInput = document.getElementById('messenger-file-input');
const uploadBtn = document.getElementById('get-stats-btn');
const cacheCheckBox = document.getElementById('cache-stats-checkbox');

const notyf = new Notyf({
    duration: 3000,
    ripple: false,
    dismissible: true,
    position: {
        x: 'center',
        y: 'top',
    },
});

async function uploadFiles() {
    const data = new FormData();
    for (const file of fileInput.files) {
        data.append('messenger_files', file, file.name);
    }

    let url = '/api/upload';
    if (cacheCheckBox.checked) {
        url += '?cache=true';
    }

    const r = await fetch(url, {
        method: 'POST',
        body: data,
    });

    const json = await r.json();
    if (json.error !== '') {
        notyf.error(json.error);
        return;
    }

    window.location.href = `/stats?id=${json.id}`;
}

window.addEventListener('load', async () => {
    uploadBtn.addEventListener('click', async () => { await uploadFiles(); });
});
