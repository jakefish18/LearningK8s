import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
    vus: 90,
    duration: '10s',
    thresholds: {
        http_req_failed: ['rate<0.01'],
    },
}

const fibonaciRange = [33, 34, 35, 36]

export default () => {
    const randomFibonaci = fibonaciRange[Math.floor(Math.random() * fibonaciRange.length)]
    const resp = http.get(`http://85.92.111.52:8080/api/v1/fibonacci/${randomFibonaci}`)
    check(resp, { '200': (r) => r.status === 200})
    sleep(1)
}