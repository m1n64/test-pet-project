import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 500,
  duration: '10m',
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8000';
const TIMEOUT = Number(__ENV.TIMEOUT || '30000');

export default function () {
  const payload = JSON.stringify({
    jsonrpc: '2.0',
    method: 'telegram.send',
    params: {
        to: '@your_channel',
        message: 'Test message from k6 smoke test',
    },
    id: Math.floor(Math.random()*1e9),
  });

  const res = http.post(`${BASE_URL}/rpc`, payload, {
    headers: { 'Content-Type': 'application/json', "X-API-KEY": "secret" },
    timeout: TIMEOUT,
  });

  check(res, {
    'status is 200': r => r.status === 200,
    'jsonrpc version ok': r => r.json('jsonrpc') === '2.0',
    'pong is true': r => r.json('result.queued') === true,
    'has id': r => r.json('id') !== undefined,
  });

  sleep(1);
}
