import { NestFactory } from "@nestjs/core";
import { SwaggerModule, DocumentBuilder } from "@nestjs/swagger";
import { ProducerServiceModule } from "./producer-service.module";

async function bootstrap() {
  const app = await NestFactory.create(ProducerServiceModule);

  const config = new DocumentBuilder()
    .setTitle("Producer Service API")
    .setDescription("The Producer Service API description")
    .setVersion("1.0")
    .build();

  const document = SwaggerModule.createDocument(app, config);
  SwaggerModule.setup("api", app, document);

  await app.listen(process.env.PRODUCER_PORT || 3000);
  console.log('🚀 Swagger running at http://localhost:' + (process.env.PRODUCER_PORT || 3000) + '/api');

}
bootstrap();
