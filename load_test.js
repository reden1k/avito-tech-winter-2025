import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  stages: [
    { duration: "10s", target: 500 },
    { duration: "10s", target: 1000 },
    { duration: "10s", target: 1500 },
    { duration: "10s", target: 2000 },
    { duration: "5s", target: 2500 },
  ]
};

export default function () {
  let res = http.post('http://localhost:8080/api/auth', JSON.stringify({
    username: 'testuser',
    password: 'password123',
  }), { headers: { 'Content-Type': 'application/json' } });

  check(res, {
    'is status 200': (r) => r.status === 200,
  });

  sleep(1); // Пауза между запросами
}
