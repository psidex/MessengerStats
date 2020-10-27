const fileInput = document.getElementById('messenger-file-input');
const uploadBtn = document.getElementById('get-stats-btn');
const infoSection = document.getElementsByClassName('website-information')[0];
const chartSection = document.getElementsByClassName('charts')[0];
const conversationTitle = document.getElementById('conversation-title');

let webSocketUrl = window.location.protocol === 'https:' ? 'wss' : 'ws';
webSocketUrl += `://${window.location.host}/api/ws`;

function setTitle(titleText) {
    conversationTitle.textContent = `${titleText}`;
}

function hideInfoShowCharts() {
    infoSection.style.display = 'none';
    chartSection.style.visibility = 'visible';
}

function hideChartsShowInfo() {
    setTitle('');
    infoSection.style.display = 'block';
    chartSection.style.visibility = 'hidden';
}

function createCharts(jsonData) {
    const messengerColour = 'rgb(0,198,255)';

    Highcharts.chart('messages-per-month-chart', {
        credits: { enabled: false },
        chart: { type: 'line' },
        title: { text: 'Messages Per Month' },
        xAxis: { categories: jsonData.messages_per_month.categories },
        yAxis: { title: { text: 'Messages Sent' } },
        series: [{ name: 'Messages Sent', data: jsonData.messages_per_month.data }],
        legend: { enabled: false },
        plotOptions: {
            series: { lineWidth: 3 },
            line: { color: messengerColour, marker: { enabled: false } },
        },
    });

    Highcharts.chart('messages-per-user-chart', {
        credits: { enabled: false },
        chart: { type: 'pie' },
        title: { text: 'Messages Per User' },
        series: [{ name: 'Messages Sent', data: jsonData.messages_per_user }],
    });

    Highcharts.chart('messages-per-weekday-chart', {
        credits: { enabled: false },
        chart: { polar: true, type: 'line' },
        title: { text: 'Messages Per Weekday' },
        xAxis: { categories: jsonData.messages_per_weekday.categories },
        yAxis: { min: 0 },
        series: [{ name: 'Messages Sent', data: jsonData.messages_per_weekday.data }],
        legend: { enabled: false },
        plotOptions: {
            series: { lineWidth: 3 },
            line: { color: messengerColour, marker: { enabled: false } },
        },
    });
}

function uploadFiles() {
    const ws = new WebSocket(webSocketUrl);

    ws.onopen = () => {
        const fileCountByte = new Uint8Array(1);
        fileCountByte[0] = fileInput.files.length;
        ws.send(fileCountByte);

        // Send off all the files asap.
        for (let i = 0; i < fileInput.files.length; i++) {
            console.time(`read file ${i} into memory`);
            const file = fileInput.files[i];
            file.arrayBuffer().then((data) => {
                console.timeEnd(`read file ${i} into memory`);
                ws.send(data);
            });
        }
    };

    ws.onmessage = async (e) => {
        ws.close();
        const data = JSON.parse(await e.data.text());

        let title = data.conversation_title;
        if (data.error !== '') {
            title = `Error: ${data.error}`;
        } else {
            // No error so load charts.
            createCharts(data);
        }

        history.pushState({ stats: true, title }, null, '/stats');
        setTitle(title);
        hideInfoShowCharts();
    };
}

window.addEventListener('load', () => {
    if (navigator.userAgent.toLowerCase().indexOf('firefox') > -1) {
        alert('This app currently does not work in Firefox due to a bug in Firefox');
        return;
    }

    uploadBtn.addEventListener('click', uploadFiles);

    window.addEventListener('popstate', (e) => {
        if (e.state === null || e.state.stats !== true) {
            hideChartsShowInfo();
        } else if (e.state.stats === true) {
            hideInfoShowCharts();
            setTitle(e.state.title);
        }
    });
});
