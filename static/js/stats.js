const messengerColour = 'rgb(0,198,255)';

function createMessagesPerMonthChart(jsonData) {
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
}

function createMessagesPerUserChart(jsonData) {
    Highcharts.chart('messages-per-user-chart', {
        credits: { enabled: false },
        chart: { type: 'pie' },
        title: { text: 'Messages Per User' },
        series: [{ name: 'Messages Sent', data: jsonData.messages_per_user }],
    });
}

function createMessagesPerWeekdayChart(jsonData) {
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

function setTitle(titleText) {
    document.getElementById('title').textContent = `${titleText}`;
}

window.addEventListener('load', async () => {
    const currentUrl = new URL(window.location.href);
    const id = currentUrl.searchParams.get('id');

    if (id === null) {
        setTitle('Error: no id supplied');
        return;
    }

    const rawData = await fetch(`/api/stats?id=${id}`);
    const data = await rawData.json();

    if (data.error !== '') {
        setTitle(`Error: ${data.error}`);
        return;
    }

    setTitle(data.conversation_title);
    createMessagesPerMonthChart(data);
    createMessagesPerUserChart(data);
    createMessagesPerWeekdayChart(data);
});
