import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 100,
      timeUnit: '1s',
      duration: '30m',
      preAllocatedVUs: 50,
      maxVUs: 1000,
    },
  },
  thresholds: {
    http_req_failed: ['rate<0.05'], // http errors should be less than 5%
    http_req_duration: ['p(95)<5000'], // 95% of requests should be below 5s
  },
};

const BASE_URL = __ENV.BASE_URL || 'http://arch.homework:8080';

export function setup() {
  const loginRes = http.post(`${BASE_URL}/auth/login`, JSON.stringify({
    email: 'supervisor@test.com',
    password: 'Supervisor123!',
    portal_id: 'default'
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  if (loginRes.status !== 200) {
    console.error(`Login failed during setup: ${loginRes.status} ${loginRes.body}`);
    return { authToken: '' };
  }

  console.log('Login successful');
  const authToken = loginRes.json('access_token');
  return { authToken };
}

export default function (data) {
  if (!data.authToken) {
    console.error('No auth token available in iteration');
    sleep(1);
    return;
  }

  const params = {
    headers: {
      'Authorization': `Bearer ${data.authToken}`,
      'Content-Type': 'application/json',
    },
  };

  // Select a random operation to simulate varied load
  const rand = Math.random();
  let res;

  if (rand < 0.4) {
    // 40% - Validate Token
    res = http.post(`${BASE_URL}/auth/validate`, JSON.stringify({
      access_token: data.authToken
    }), params);
  } else if (rand < 0.7) {
    // 30% - List Courses
    res = http.get(`${BASE_URL}/courses?portalId=default`, params);
  } else if (rand < 0.9) {
    // 20% - Get Portal
    res = http.get(`${BASE_URL}/portals/default`, params);
  } else {
    // 10% - List Portals
    res = http.get(`${BASE_URL}/portals`, params);
  }

  const success = check(res, {
    'is status 200': (r) => r.status === 200,
  });

  if (!success) {
    console.error(`Request failed: ${res.method} ${res.url} - status: ${res.status} body: ${res.body}`);
  }
}
