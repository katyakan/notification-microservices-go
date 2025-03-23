import { Module } from '@nestjs/common';
import { ProducerServiceController } from './producer-service.controller';
import { ProducerServiceService } from './producer-service.service';

@Module({
  imports: [],
  controllers: [ProducerServiceController],
  providers: [ProducerServiceService],
})
export class ProducerServiceModule {}
