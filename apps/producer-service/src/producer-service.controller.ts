import { Controller, Get } from '@nestjs/common';
import { ProducerServiceService } from './producer-service.service';

@Controller()
export class ProducerServiceController {
  constructor(private readonly producerServiceService: ProducerServiceService) {}

  @Get()
  getHello(): string {
    return this.producerServiceService.getHello();
  }
}
