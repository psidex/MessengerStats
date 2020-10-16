const messengerColour = 'rgb(0,198,255)';

function createMessagesPerMonthChart(jsonData) {
    const data = [];
    const categories = [];

    for (const year in jsonData.messages_per_month) {
        for (const month in jsonData.messages_per_month[year]) {
            categories.push(`${year}-${month}`);
            data.push(jsonData.messages_per_month[year][month]);
        }
    }

    Highcharts.chart('messages-per-month-chart', {
        chart: { type: 'line' },
        title: { text: 'Messages Per Month' },
        xAxis: { categories },
        yAxis: { title: { text: 'Messages Sent' } },
        series: [{ name: 'Messages Sent', data }],
        legend: { enabled: false },
        plotOptions: {
            series: { lineWidth: 3 },
            line: { color: messengerColour, marker: { enabled: false } },
        },
    });
}

function createMessagesPerUserChart(jsonData) {
    const data = [];

    for (const user in jsonData.messages_per_user) {
        data.push({
            name: user,
            y: jsonData.messages_per_user[user],
        });
    }

    Highcharts.chart('messages-per-user-chart', {
        chart: { type: 'pie' },
        title: { text: 'Messages Per User' },
        series: [{ name: 'Messages Sent', data }],
    });
}

function createMessagesPerWeekdayChart(jsonData) {
    const data = [];
    const categories = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];

    categories.forEach((weekday) => {
        data.push(jsonData.messages_per_weekday[weekday]);
    });

    Highcharts.chart('messages-per-weekday-chart', {
        chart: { polar: true, type: 'line' },
        title: { text: 'Messages Per Weekday' },
        xAxis: { categories },
        yAxis: { min: 0 },
        series: [{ name: 'Messages Sent', data }],
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
        return;
    }

    const rawData = await fetch(`/api/stats?id=${id}`);
    const data = await rawData.json();

    if (data.error === 'ID not found') {
        alert('ID not found');
        return;
    }

    setTitle(data.conversation_title);
    createMessagesPerMonthChart(data);
    createMessagesPerUserChart(data);
    createMessagesPerWeekdayChart(data);
});
