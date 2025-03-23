import { Module } from "@nestjs/common";
import { ProducerServiceController } from "./producer-service.controller";
import { ProducerServiceService } from "./producer-service.service";
import { ProducerModule } from "./producer.module";

@Module({
  imports: [ProducerModule],
  controllers: [ProducerServiceController],
  providers: [ProducerServiceService],
})
export class ProducerServiceModule {}
