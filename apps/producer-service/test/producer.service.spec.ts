import { Test, TestingModule } from "@nestjs/testing";
import { KafkaService } from "../src/services/kafka.service";

describe("KafkaService", () => {
  let service: KafkaService;

  beforeEach(async () => {
    const module: TestingModule = await Test.createTestingModule({
      providers: [KafkaService],
    }).compile();

    service = module.get<KafkaService>(KafkaService);
  });

  it("should be defined", () => {
    expect(service).toBeDefined();
  });

  it("should send message successfully", async () => {
    const message = {
      type: "test",
      payload: { text: "test message" },
    };

    const result = await service.sendMessage(message);
    expect(result).toHaveProperty("id");
  });
});
