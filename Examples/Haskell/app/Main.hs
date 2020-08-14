{-# LANGUAGE OverloadedStrings #-}
module Main where

import qualified Data.ByteString.Char8 as BC
import Network.HTTP.Simple
import Network.HTTP.Client.Types

host :: BC.ByteString
host = "localhost"

port :: Int
port = 8080

path :: BC.ByteString
path = "/cache/Hello"

value = "World"

request :: BC.ByteString -> Request
request method = setRequestMethod method
                                  $ setRequestHost host
                                  $ setRequestPath path
                                  $ setRequestBody value
                                  $ setRequestPort port
                                  $ defaultRequest

main :: IO ()
main = do
  httpNoBody $ request "PUT"
  response <- httpBS $ request "GET"
  let status = getResponseStatusCode response
  let value = getResponseBody response
  print $ "Hello = " <> value

