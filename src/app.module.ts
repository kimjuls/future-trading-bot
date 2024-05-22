import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { CandleCollectorModule } from './candle-collector/candle-collector.module';
import { StrategyModule } from './strategy/strategy.module';
import { CandleAnalyzerModule } from './candle-analyzer/candle-analyzer.module';
import { OrderExecutorModule } from './order-executor/order-executor.module';
import { MyLoggerModule } from './my-logger/my-logger.module';
import { RiskManagerModule } from './risk-manager/risk-manager.module';

@Module({
  imports: [
    CandleCollectorModule,
    StrategyModule,
    CandleAnalyzerModule,
    OrderExecutorModule,
    MyLoggerModule,
    RiskManagerModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
