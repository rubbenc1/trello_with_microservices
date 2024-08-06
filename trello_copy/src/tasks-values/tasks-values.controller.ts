import { Body, Controller, Delete, Get, HttpException, HttpStatus, Param, ParseIntPipe, Post, Put, UseGuards } from '@nestjs/common';
import { ApiOperation, ApiResponse, ApiTags } from '@nestjs/swagger';
import { TasksValuesService } from './tasks-values.service';
import { TasksValuesDto } from './dto/tasks-values.dto';
import { ColumnOwnerGuard } from 'src/auth/column-owner.guard';


@ApiTags('Tasks column values')
@Controller('tasks-values/tasks/:taskId/columns/:columnId/values')
export class TasksValuesController {

    constructor(
        private tasksValuesService: TasksValuesService
    ){}

    @ApiOperation({summary: 'Create value for the column'})
    @ApiResponse({status: 201, description: 'Value created successfully'})
    @Post()
    @UseGuards(ColumnOwnerGuard)
    async createColumnValue(
      @Param('taskId', ParseIntPipe) taskId: number,
      @Param('columnId', ParseIntPipe) columnId: number,
      @Body() value: TasksValuesDto
    ) {

      try {
        const response = await this.tasksValuesService.createTaskValues(taskId, columnId,value);
        return { status: 'Message sent to queue', response };
      } catch (error) {
        console.error('Error sending message to queue:', error);
        throw new HttpException('Failed to send message', HttpStatus.INTERNAL_SERVER_ERROR);
      }
    }

    @ApiOperation({summary: 'Update value for the column'})
    @ApiResponse({status: 200, description: 'Value updated successfully'})
    @UseGuards(ColumnOwnerGuard)
    @Put(':valueId')
    async updateColumnValue(
      @Param('valueId', ParseIntPipe) valueId: number,
      @Body() value: TasksValuesDto
    ){
      const message = {
        action: "update",
        valueId,
        ...value,
      };

      try {
        const response = await this.tasksValuesService.updateTaskValues(valueId,value);
        return { status: 'Message sent to queue', response };
      } catch (error) {
        console.error('Error sending message to queue:', error);
        throw new HttpException('Failed to send message', HttpStatus.INTERNAL_SERVER_ERROR);
      }
    }

    @ApiOperation({summary: 'Delete value for the column'})
    @ApiResponse({status: 200, description: 'Value deleted successfully'})
    @UseGuards(ColumnOwnerGuard)
    @Delete(':valueId')
    async deleteColumnValue(@Param('valueId', ParseIntPipe) valueId: number){
      const message = {
        action: "delete",
        valueId
      }
      try {
        const response = await this.tasksValuesService.deleteTaskValues(valueId);
        return { status: 'Message sent to queue', response };
      } catch (error) {
        console.error('Error sending message to queue:', error);
        throw new HttpException('Failed to send message', HttpStatus.INTERNAL_SERVER_ERROR);
      }
    }

    @ApiOperation({summary: 'Get value for the column'})
    @ApiResponse({status: 200, description: 'Value retrieved successfully'})
    @ApiResponse({status: 500, description: 'Failed to retrieve value'})
    @UseGuards(ColumnOwnerGuard)
    @Get(':valueId')
    async getColumnValue(
      @Param('valueId', ParseIntPipe) valueId: number
    ) {
      const message = {
        action: "get",
        valueId
      }
      try {
        const response = await this.tasksValuesService.getTaskValues(valueId);
        return { status: 'Message sent to queue', response };
      } catch (error) {
        console.error('Error sending message to queue:', error);
        throw new HttpException('Failed to send message', HttpStatus.INTERNAL_SERVER_ERROR);
      }
    }
}

