import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  stages: [
    { duration: '10s', target: 10 },   // разогрев: 10 пользователей
    { duration: '30s', target: 50 },   // нагрузка: 50 пользователей
  ],
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
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post(url, payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(1);
}
