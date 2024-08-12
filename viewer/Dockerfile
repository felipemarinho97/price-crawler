# Use the official Node.js Alpine base image for smaller size
FROM node:18-alpine AS build

# Set working directory
WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install --frozen-lockfile

# Copy the rest of the application
COPY . .

# Inject environment variable 
ARG DATA_BUCKET_URL
ENV REACT_APP_DATA_BUCKET_URL=$DATA_BUCKET_URL

# Build the React app
RUN npm run build

# Use a lightweight web server to serve the built files
FROM nginx:alpine

# Copy the build output to the Nginx html directory
COPY --from=build /app/build /usr/share/nginx/html

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]