import { Injectable } from '@nestjs/common';

@Injectable()
export class ProducerServiceService {
  getHello(): string {
    return 'Hello World!';
  }
}
