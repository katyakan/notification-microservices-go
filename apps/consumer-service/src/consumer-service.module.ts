import { Module } from '@nestjs/common';
import { ConsumerServiceController } from './consumer-service.controller';
import { ConsumerServiceService } from './consumer-service.service';

@Module({
  imports: [],
  controllers: [ConsumerServiceController],
  providers: [ConsumerServiceService],
})
export class ConsumerServiceModule {}
