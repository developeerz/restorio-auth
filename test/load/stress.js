import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    high_rps_test: {
      executor: 'constant-arrival-rate',
      rate: 1000,          // 1000 RPS
      timeUnit: '1s',
      duration: '30s',     // Продолжительность теста
      preAllocatedVUs: 100,
      maxVUs: 1000,
    },
  },
};

function randomUser() {
    const randomInt = Math.floor(Math.random() * 100000);
    return {
      firstname: `firstname${randomInt}`,
      lastname: `lastname${randomInt}`,
      telegram: `telegram${randomInt}`,
      password: `pass${randomInt}`,
    };
  }

export default function () {
  const url = 'http://localhost/api/web-gateway/user/sign-up';
  const payload = JSON.stringify(randomUser());

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'status is 2xx': (r) => r.status >= 200 && r.status < 300,
    'response time OK': (r) => r.timings.duration < 100,
  });
}
