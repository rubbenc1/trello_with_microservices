import { Inject, Injectable, NotFoundException } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import { InjectRepository } from '@nestjs/typeorm';
import { lastValueFrom } from 'rxjs';
import { TasksColumnOptions } from 'src/tasks-column-options/tasks-column-options.model';
import { TasksColumn } from 'src/tasks-column/tasks-column.model';
import { Tasks } from 'src/tasks/tasks.model';
import { Repository } from 'typeorm';
import { TasksValuesDto } from './dto/tasks-values.dto';

@Injectable()
export class TasksValuesService {
  constructor(
    @Inject('TASKS_VALUES_SERVICE') private client: ClientProxy,
    @InjectRepository(Tasks) private tasksRepository: Repository<Tasks>,
    @InjectRepository(TasksColumn) private tasksColumnsRepository: Repository<TasksColumn>,
    @InjectRepository(TasksColumnOptions) private tasksColumnOptionsRepository: Repository<TasksColumnOptions>
  ) {}

  private setValue(value: string | number) {
    let stringValue: string | null = null;
    let numberValue: number | null = null;
    let type: 'string' | 'number';

    if (typeof value === 'string') {
      type = 'string';
      stringValue = value;
    } else if (typeof value === 'number') {
      type = 'number';
      numberValue = value;
    }

    return { stringValue, numberValue, type };
  }

  async createTaskValues(taskId: number, columnId: number, valueDto: TasksValuesDto) {
    // Verify the task exists
    const task = await this.tasksRepository.findOne({ where: { id: taskId } });
    if (!task) {
      throw new NotFoundException(`Task with ID ${taskId} not found`);
    }

    // Verify the column exists
    const column = await this.tasksColumnsRepository.findOne({ where: {columnId } });
    if (!column) {
      throw new NotFoundException(`Column with ID ${columnId} not found`);
    }

    const { value, optionId } = valueDto;
    let valueData = { stringValue: null, numberValue: null, type: null };

    if (optionId) {
      const option = await this.tasksColumnOptionsRepository.findOne({ where: {optionId } });
      if (!option) {
        throw new NotFoundException(`Option with ID ${optionId} not found`);
      }
      // When optionId is present, it indicates the use of a predefined option
      valueData = { stringValue: null, numberValue: null, type: 'number' };
    } else {
      valueData = this.setValue(value);
    }

    const message = {
      action: 'create',
      taskId,
      columnId,
      optionId,
      ...valueData,
    };

    try {
      // Send the JSON-formatted message to RabbitMQ
      await lastValueFrom(this.client.emit('tasks_queue', message));
    } catch (error) {
      console.error('Error while emitting to queue:', error);
      throw error;
    }
  }

  async updateTaskValues(valueId: number, valueDto: TasksValuesDto) {
    const { value, optionId } = valueDto;
    let valueData = { stringValue: null, numberValue: null, type: null };

    if (optionId) {
      const option = await this.tasksColumnOptionsRepository.findOne({ where: {optionId } });
      if (!option) {
        throw new NotFoundException(`Option with ID ${optionId} not found`);
      }
      valueData = { stringValue: null, numberValue: null, type: 'number' };
    } else {
      valueData = this.setValue(value);
    }

    const message = {
      action: 'update',
      valueId,
      optionId,
      ...valueData,
    };

    try {
      // Send the JSON-formatted message to RabbitMQ
      await lastValueFrom(this.client.emit('tasks_queue', message));
    } catch (error) {
      console.error('Error while emitting to queue:', error);
      throw error;
    }
  }

  async deleteTaskValues(valueId: number) {
    const message = {
      action: 'delete',
      valueId,
    };

    try {
      // Send the JSON-formatted message to RabbitMQ
      await lastValueFrom(this.client.emit('tasks_queue', message));
    } catch (error) {
      console.error('Error while emitting to queue:', error);
      throw error;
    }
  }

  async getTaskValues(valueId: number) {
    const message = {
      action: 'get',
      valueId,
    };

    try {
      // Send the JSON-formatted message to RabbitMQ and await the response
      return await lastValueFrom(this.client.send('tasks_queue', message));
    } catch (error) {
      console.error('Error while sending get request to queue:', error);
      throw error;
    }
  }

  async deleteTaskOfAllValues(id: number) {
    const taskToDelete = await this.tasksRepository.findOne({ where: { id } });
    if (!taskToDelete) {
        throw new NotFoundException(`Task with ID ${id} not found`);
        }
    const message = {
        action: 'deleteAllForTask',
        taskId: id,
    };
    try {
        await lastValueFrom(this.client.emit('tasks_queue', message));
      } catch (error) {
        console.error('Error while emitting to queue:', error);
        throw error;
      }
    
    }
}

