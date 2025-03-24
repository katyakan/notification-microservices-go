import { Module } from "@nestjs/common";
import { ConsumerController } from "./controllers/consumer.controller";
import { KafkaService } from "./services/kafka.service";
import { ConfigModule } from '@nestjs/config';
@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: '.env',
    }),
  ],
  controllers: [ConsumerController],
  providers: [KafkaService],
})
export class ConsumerModule {}
