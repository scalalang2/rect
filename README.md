# RECT
RECT: Reducing cross-shard transactions for efficient load balancing in ethereum sharding environment.

## Motivation
크로스-샤드 트랜잭션은 싱글-샤드 트랜잭션 대비 3x 높은 연산력을 요구한다. 샤딩 환경에서 빈번한 크로스-샤드 트랜잭션은 네트워크의 성능 저하로 이어지기 때문에 최대한 크로스-샤드 트랜잭션이 적게 발생하도록 트랜잭션을 각 샤드에 할당할 필요가 있다.
