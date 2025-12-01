const stompClient = new StompJs.Client({
  reconnectDelay: 5000,
  webSocketFactory: () => new SockJS('/ws'),
});

const quotes = new Map();

stompClient.onConnect = frame => {
  setConnected(true);
  console.log('Connected: ' + frame);
  // stompClient.subscribe('/topic/quotes', message => showMessage(message.body));
  stompClient.subscribe('/topic/quotes', message => updateQuote(message.body));
};

stompClient.onWebSocketError = error => {
  console.error('WebSocket error', error);
};

stompClient.onStompError = frame => {
  console.error('Broker reported STOMP error: ' + frame.headers['message']);
  console.error('Additional details: ' + frame.body);
};
function setConnected(connected) {
  $('#connect').prop('disabled', connected);
  $('#disconnect').prop('disabled', !connected);
  $('#quotesTable').css('display', connected ? 'table' : 'none');
  if (!connected) {
    quotes.clear();
    $('#quotesBody').empty();
  }
}
// function setConnected(connected) {
//   $('#connect').prop('disabled', connected);
//   $('#disconnect').prop('disabled', !connected);
//   if (connected) {
//     $('#conversation').show();
//   } else {
//     $('#conversation').hide().empty();
//   }
// }

function connect() {
  stompClient.activate();
}

function disconnect() {
  stompClient.deactivate();
  setConnected(false);
  console.log('Disconnected');
}

// function showMessage(payload) {
//   let text = payload;
//   try {
//     const data = JSON.parse(payload);
//     text = JSON.stringify(data, null, 2);
//   } catch (ignore) {
//     // payload was not JSON; keep raw text
//   }

//   const $line = $('<pre/>').text(text);
//   $('#conversation').append($line);
// }

function updateQuote(payload) {
  try {
    const data = JSON.parse(payload);
    const symbol = data.symbol ?? 'N/A';
    const price = Number(data.lastPrice ?? 0);
    const change = Number(data.change ?? 0);
    // const changePct = data.lastPrice && data.change
    //   ? (change / (price - change)) * 100
    //   : 0;
    let changePct = 0;
    if (!Number.isNaN(change) && price != null) {
      const previous = price - change;
      if (previous > 0) {
        changePct = (change / previous) * 100;
      } else if (price !== 0) {
        // 後端也許會直接給 changePct：優先使用
        changePct = data.changePct ?? (change / price) * 100;
      }
    }    

    let updatedAt = '';
    if (data.tradeTime) {
      updatedAt = data.tradeTime;
    } else if (data.tradeTs) {
      updatedAt = new Date(data.tradeTs)
        .toLocaleTimeString('zh-TW', { hour12: false });
    } else if (data.ts) {
      updatedAt = new Date(data.ts)
        .toLocaleTimeString('zh-TW', { hour12: false });
    }


    quotes.set(symbol, { price, change, changePct, updatedAt });
    renderQuotes();
  } catch (err) {
    console.warn('無法解析行情 payload：', payload, err);
  }
}


function renderQuotes() {
  const tbody = $('#quotesBody').empty();
  for (const [symbol, { price, change, changePct, updatedAt }] of quotes.entries()) {
    const cls = change > 0 ? 'up' : change < 0 ? 'down' : 'neutral';
    tbody.append(`
      <tr>
        <td>${symbol}</td>
        <td>${price.toFixed(2)}</td>
        <td class="${cls}">${change.toFixed(2)}</td>
        <td class="${cls}">${changePct.toFixed(2)}%</td>
        <td>${updatedAt}</td>
      </tr>
    `);
  }
}

$(function () {
  $('form').on('submit', event => event.preventDefault());
  $('#connect').click(connect);
  $('#disconnect').click(disconnect);

  // 初始化狀態
  // setConnected(false);
});

