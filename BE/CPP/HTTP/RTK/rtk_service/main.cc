#include <drogon/drogon.h>
#include <iostream>
int main() {
    //Set HTTP listener address and port
    drogon::app().addListener("0.0.0.0",8001);
    //Load config file
    //drogon::app().loadConfigFile("../config.json");
    //Run HTTP framework,the method will block in the internal event loop
    drogon::app().setLogLevel(trantor::Logger::LogLevel::kDebug);
    std::cout << "start server at 8001...." << std::endl;
    drogon::app().run();
    return 0;
}
