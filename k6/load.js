import http from 'k6/http';
import { check, sleep } from 'k6';

const baseURL = __ENV.BASE_URL || 'http://localhost:8080';

export const options = {
  vus: 5,
  duration: '30s',
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<500'],
  },
};

export default function () {
  const response = http.get(`${baseURL}/api/hello?name=k6`);
  check(response, {
    'status is 200': (res) => res.status === 200,
    'response has greeting': (res) => res.json('message') === 'hello, k6',
  });
  sleep(1);
}
