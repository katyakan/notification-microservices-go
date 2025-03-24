import { Module } from "@nestjs/common";
import { NotificationServiceController } from "./notification-service.controller";
import { NotificationServiceService } from "./notification-service.service";
import { KafkaService } from "./services/kafka.service";
import { TelegramService } from "./services/telegram.service";
import { ConfigModule } from '@nestjs/config';
@Module({
  imports: [
    ConfigModule.forRoot({
      isGlobal: true,
      envFilePath: '.env',
    }),
  ],
  controllers: [NotificationServiceController],
  providers: [NotificationServiceService, KafkaService, TelegramService],
})
export class NotificationServiceModule {}
