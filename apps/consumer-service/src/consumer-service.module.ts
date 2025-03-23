import { Module } from "@nestjs/common";
import { ConsumerServiceController } from "./consumer-service.controller";
import { ConsumerServiceService } from "./consumer-service.service";
import { ConsumerModule } from "./consumer.module";

@Module({
  imports: [ConsumerModule],
  controllers: [ConsumerServiceController],
  providers: [ConsumerServiceService],
})
export class ConsumerServiceModule {}
