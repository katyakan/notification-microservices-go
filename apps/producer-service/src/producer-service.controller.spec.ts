import { Test, TestingModule } from '@nestjs/testing';
import { ProducerServiceController } from './producer-service.controller';
import { ProducerServiceService } from './producer-service.service';

describe('ProducerServiceController', () => {
  let producerServiceController: ProducerServiceController;

  beforeEach(async () => {
    const app: TestingModule = await Test.createTestingModule({
      controllers: [ProducerServiceController],
      providers: [ProducerServiceService],
    }).compile();

    producerServiceController = app.get<ProducerServiceController>(ProducerServiceController);
  });

  describe('root', () => {
    it('should return "Hello World!"', () => {
      expect(producerServiceController.getHello()).toBe('Hello World!');
    });
  });
});
