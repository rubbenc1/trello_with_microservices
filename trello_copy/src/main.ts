import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import * as dotenv from 'dotenv';
import { DocumentBuilder, SwaggerModule } from '@nestjs/swagger';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';

async function start() {
    dotenv.config();

    const PORT = process.env.PORT || 3000;
    const app = await NestFactory.create(AppModule);

    const config = new DocumentBuilder()
        .setTitle('Project Trello')
        .setDescription('API Documentation')
        .setVersion('1.0.0')
        .addTag('rubbenc')
        .build();
    const document = SwaggerModule.createDocument(app, config);
    SwaggerModule.setup('/api/docs', app, document);


    await app.listen(PORT, () => {
        console.log(`Server started on port = ${PORT}`);
    });
}

start();
