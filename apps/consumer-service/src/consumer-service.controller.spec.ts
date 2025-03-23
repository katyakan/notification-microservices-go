import { Test, TestingModule } from "@nestjs/testing";
import { ConsumerServiceController } from "./consumer-service.controller";
import { ConsumerServiceService } from "./consumer-service.service";

describe("ConsumerServiceController", () => {
  let consumerServiceController: ConsumerServiceController;

  beforeEach(async () => {
    const app: TestingModule = await Test.createTestingModule({
      controllers: [ConsumerServiceController],
      providers: [ConsumerServiceService],
    }).compile();

    consumerServiceController = app.get<ConsumerServiceController>(ConsumerServiceController);
  });

  describe("root", () => {
    it("should return "Hello World!"", () => {
      expect(consumerServiceController.getHello()).toBe("Hello World!");
    });
  });
});
