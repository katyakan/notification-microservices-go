import { Injectable, OnModuleInit } from "@nestjs/common";
import { Telegraf } from "telegraf";

@Injectable()
export class TelegramService implements OnModuleInit {
  private bot: Telegraf;

  constructor() {
    if (!process.env.TELEGRAM_BOT_TOKEN) {
      throw new Error("TELEGRAM_BOT_TOKEN is not set");
    }
    this.bot = new Telegraf(process.env.TELEGRAM_BOT_TOKEN);
  }

  async onModuleInit() {
    await this.bot.launch();
    console.log("Telegram bot started");
  }

  async sendMessage(chatId: number, text: string): Promise<void> {
    try {
      await this.bot.telegram.sendMessage(chatId, text);
      console.log(`Message sent to chat ${chatId}`);
    } catch (error) {
      console.error("Error sending message to Telegram:", error);
      throw error;
    }
  }
}
