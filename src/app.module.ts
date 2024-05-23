import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { CandleCollectorModule } from './candle-collector/candle-collector.module';
import { StrategyModule } from './strategy/strategy.module';
import { CandleAnalyzerModule } from './candle-analyzer/candle-analyzer.module';
import { MyLoggerModule } from './my-logger/my-logger.module';
import { RiskManagerModule } from './risk-manager/risk-manager.module';
import { OrderManagerModule } from './order-manager/order-manager.module';
import { BinanceApiService } from './binance-api/binance-api.service';
import { BinanceApiModule } from './binance-api/binance-api.module';

@Module({
  imports: [
    CandleCollectorModule,
    StrategyModule,
    CandleAnalyzerModule,
    MyLoggerModule,
    RiskManagerModule,
    OrderManagerModule,
    BinanceApiModule,
  ],
  controllers: [AppController],
  providers: [AppService, BinanceApiService],
})
export class AppModule {}
