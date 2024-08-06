import { Module } from '@nestjs/common';
import { TasksValuesController } from './tasks-values.controller';
import { TasksValuesService } from './tasks-values.service';
import { AuthModule } from 'src/auth/auth.module';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { TypeOrmModule } from '@nestjs/typeorm';
import { TasksColumn } from 'src/tasks-column/tasks-column.model';
import { TasksColumnOptions } from 'src/tasks-column-options/tasks-column-options.model';
import { Tasks } from 'src/tasks/tasks.model';
import { UserProjects } from 'src/projects/user-projects.model';


@Module({
  controllers: [TasksValuesController],
  providers: [TasksValuesService],
  imports: [
    AuthModule,
    TypeOrmModule.forFeature([Tasks,TasksColumn, TasksColumnOptions, UserProjects]),
    ClientsModule.register([
    {
      name: "TASKS_VALUES_SERVICE",
      transport: Transport.RMQ, 
      options: {
        urls: [`amqp://localhost:${process.env.URI_PORT || 5672}`],
        queue: 'tasks_queue',
        queueOptions:{
          durable: true
        },
      },
    },
  ]),
  ], 
  exports: [TasksValuesService]
})
export class TasksValuesModule {}
