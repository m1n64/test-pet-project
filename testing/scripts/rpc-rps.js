import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    rpc_rate: {
      executor: 'constant-arrival-rate',
      rate: 200,
      timeUnit: '1s',
      duration: '5m',
      preAllocatedVUs: 100, 
      maxVUs: 500,
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<200', 'p(99)<500'],
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8000';
const TIMEOUT = Number(__ENV.TIMEOUT || '30000');

const METHODS = ['system.ping'];

export default function () {
  const method = METHODS[Math.floor(Math.random()*METHODS.length)];
  const payload = JSON.stringify({
    jsonrpc: '2.0',
    method,
    params: {},
    id: Math.floor(Math.random()*1e9),
  });

  const res = http.post(`${BASE_URL}/rpc`, payload, {
    headers: { 'Content-Type': 'application/json', "X-API-KEY": "secret" },
    timeout: TIMEOUT,
  });

  check(res, {
    'status is 200': r => r.status === 200,
    'jsonrpc version ok': r => r.json('jsonrpc') === '2.0',
    'pong is true': r => r.json('result.pong') === true,
    'has id': r => r.json('id') !== undefined,
  });
}
