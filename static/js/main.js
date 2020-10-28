//
// Constants
//

const fileInput = document.getElementById('messenger-file-input');
const uploadBtn = document.getElementById('upload-files-btn');
const uploadProgressBar = document.getElementById('upload-progress-bar');
const uploadErrorText = document.getElementById('upload-error-text');
const infoSection = document.getElementsByClassName('website-information')[0];
const chartSection = document.getElementsByClassName('charts')[0];
const conversationTitle = document.getElementById('conversation-title');

let webSocketUrl = window.location.protocol === 'https:' ? 'wss' : 'ws';
webSocketUrl += `://${window.location.host}/api/ws`;

//
// Element value setters
//

function setTitle(titleText) {
    conversationTitle.textContent = `${titleText}`;
}

function setUploadPercent(percent) {
    uploadProgressBar.value = percent;
}

function setUploadErrorText(errorText) {
    uploadErrorText.textContent = errorText;
}

//
// State setters
// There is some non-DRY code here but it's so that everything that is happening is clear.
//

function setChartViewState(title) {
    setTitle(title);
    uploadProgressBar.style.visibility = 'hidden';
    uploadErrorText.style.visibility = 'hidden';
    infoSection.style.display = 'none';
    chartSection.style.visibility = 'visible';
}

function setUploadingState() {
    setUploadErrorText('');
    setUploadPercent(0);
    uploadProgressBar.style.visibility = 'visible';
    uploadErrorText.style.visibility = 'hidden';
}

function setUploadErrorState(errorText) {
    setUploadErrorText(errorText);
    uploadProgressBar.style.visibility = 'hidden';
    uploadErrorText.style.visibility = 'visible';
}

function setInfoViewState() {
    setTitle('');
    uploadProgressBar.style.visibility = 'hidden';
    uploadErrorText.style.visibility = 'hidden';
    infoSection.style.display = 'block';
    chartSection.style.visibility = 'hidden';
}

//
// Functions & Event listeners
//

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
    setUploadingState();

    const ws = new WebSocket(webSocketUrl);

    ws.onopen = () => {
        const fileCountByte = new Uint8Array(1);
        fileCountByte[0] = fileInput.files.length;
        ws.send(fileCountByte);

        // Send off all the files asap.
        for (let i = 0; i < fileInput.files.length; i++) {
            const file = fileInput.files[i];
            file.arrayBuffer().then((data) => { ws.send(data); });
        }
    };

    ws.onmessage = async (e) => {
        const data = JSON.parse(await e.data.text());

        if (data.progress !== undefined) {
            setUploadPercent(data.progress);
        } else if (data.error !== '') {
            setUploadErrorState(`Error: ${data.error}`);
        } else {
            ws.close();
            const title = data.conversation_title;
            createCharts(data);
            history.pushState({ stats: true, title }, null, '/stats');
            setChartViewState(title);
        }
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
            setInfoViewState();
        } else if (e.state.stats === true) {
            setChartViewState(e.state.title);
        }
    });
});
