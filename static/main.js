// Note about canvas sizes:
// https://www.chartjs.org/docs/latest/general/responsive.html#important-note

const rgbColour = 'rgb(0,198,255)';

function createMessagesPerMonthChart(jsonData) {
    const data = {
        labels: [],
        datasets: [{
            data: [],
            borderColor: rgbColour,
            lineTension: 0.25,
        }],
    };

    for (const year in jsonData.messages_per_month) {
        for (const month in jsonData.messages_per_month[year]) {
            data.labels.push(`${year}-${month}`);
            data.datasets[0].data.push(jsonData.messages_per_month[year][month]);
        }
    }

    const ctx = document.querySelector('#messagesPerMonthChart').getContext('2d');
    new Chart(ctx, {
        type: 'line',
        data,
        options: {
            title: {
                display: true,
                text: 'Messages Sent Per Month',
            },
            legend: {
                display: false,
            },
            maintainAspectRatio: false,
        },
    });
}

function createMessagesPerUserChart(jsonData) {
    const data = {
        labels: [],
        datasets: [{
            data: [],
            backgroundColor: [],
            borderColor: rgbColour,
        }],
    };

    for (const user in jsonData.messages_per_user) {
        data.labels.push(user);
        data.datasets[0].data.push(jsonData.messages_per_user[user]);
    }

    const ctx = document.querySelector('#messagesPerUserChart').getContext('2d');
    new Chart(ctx, {
        type: 'pie',
        data,
        options: {
            title: {
                display: true,
                text: 'Messages Sent Per User',
            },
            legend: {
                display: false,
            },
            maintainAspectRatio: false,
            plugins: {
                labels: {
                    render: 'label',
                },
            },
        },
    });
}

function createMessagesPerWeekdayChart(jsonData) {
    const data = {
        labels: [],
        datasets: [{
            data: [],
            borderColor: rgbColour,
        }],
    };

    const weekdays = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'];
    weekdays.forEach((weekday) => {
        data.labels.push(weekday);
        data.datasets[0].data.push(jsonData.messages_per_weekday[weekday]);
    });

    const ctx = document.querySelector('#messagesPerWeekdayChart').getContext('2d');
    new Chart(ctx, {
        type: 'radar',
        data,
        options: {
            title: {
                display: true,
                text: 'Total Messages Sent Per Weekday',
            },
            legend: {
                display: false,
            },
            maintainAspectRatio: false,
            scale: {
                ticks: {
                    suggestedMin: 0,
                },
            },
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
