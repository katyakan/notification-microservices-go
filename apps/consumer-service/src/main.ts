import { NestFactory } from "@nestjs/core";
import { ConsumerModule } from "./consumer.module";


async function bootstrap() {
  const app = await NestFactory.create(ConsumerModule);
  await app.listen(process.env.CONSUMER_PORT || 3001);
  console.log('⚙️ PORT:', process.env.CONSUMER_PORT);

}
bootstrap();
