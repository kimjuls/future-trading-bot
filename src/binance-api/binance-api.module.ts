import { HttpModule } from '@nestjs/axios';
import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { BinanceApiService } from './binance-api.service';

@Module({
  imports: [ConfigModule.forRoot(), HttpModule],
  providers: [BinanceApiService],
})
export class BinanceApiModule {}
